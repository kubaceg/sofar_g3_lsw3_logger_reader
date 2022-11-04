package influx

import (
	"fmt"
	"log"
	"time"

	"git.xelasys.ro/sigxcpu/sofar/ports"
	influxdb "github.com/influxdata/influxdb1-client/v2"
)

type dbRecord struct {
	measurement map[string]interface{}
	timestamp   time.Time
}
type db struct {
	databaseURL  string
	databaseName string
	tags         map[string]string
	buffer       chan (dbRecord)
	chanStop     chan (struct{})
}

func New(databaseURL string, databaseName string, tags map[string]string) ports.Database {
	if tags == nil {
		tags = map[string]string{}
	}

	res := &db{
		databaseURL:  databaseURL,
		databaseName: databaseName,
		tags:         tags,
		buffer:       make(chan dbRecord, 1500),
	}

	go res.runLoop()
	return res
}

func (d *db) InsertRecord(measurement map[string]interface{}) error {
	measurementCopy := make(map[string]interface{}, len(measurement))
	for k, v := range measurement {
		measurementCopy[k] = v
	}
	select {
	case d.buffer <- dbRecord{measurement: measurementCopy, timestamp: time.Now()}:
		return nil
	default:
		return fmt.Errorf("buffer is full")
	}
}

func (d *db) insertRecord(r dbRecord) error {
	c, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr: d.databaseURL,
	})

	if err != nil {
		return err
	}

	defer c.Close()

	bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  d.databaseName,
		Precision: "s",
	})

	pt, err := influxdb.NewPoint("sofar", d.tags, r.measurement, r.timestamp)

	if err != nil {
		return err
	}

	bp.AddPoint(pt)

	err = c.Write(bp)

	if err != nil {
		return err
	}

	// 	T 2020/02/11 22:46:09.440463 192.168.27.140:58188 -> 109.96.190.26:80 [AP] #4
	// POST /write?db=solar&precision=s HTTP/1.1.
	// content-type: application/x-www-form-urlencoded.
	// host: influx.alex.xelasys.ro.
	// authorization: Basic Og==.
	// content-length: 651.
	// Connection: close.
	// .
	// solar,device=mpi,query=query_power_status solar_input_power_1=0,solar_input_power_2=0,battery_power=0,ac_input_active_power_r=512,ac_input_active_power_s=535,ac_input_active_power_t=93,ac_input_total_active_power=1140,ac_output_active_power_r=517,ac_output_active_power_s=502,ac_output_active_power_t=58,ac_output_total_active_power=1077,ac_output_apperent_power_r=569,ac_output_apperent_power_s=547,ac_output_apperent_power_t=73,ac_output_total_apperent_power=1189,ac_output_power_percentage=17,ac_output_connect_status=1,solar_input_1_work_status=0,solar_input_2_work_status=0,battery_power_direction=1,dc_ac_power_direction=1,line_power_direction=1

	return nil

}

func (d *db) runLoop() {
EXIT:
	for {
		select {
		case <-d.chanStop:
			log.Printf("exiting influx send loop with %d events in the queue", len(d.buffer))

			goto EXIT
		case r := <-d.buffer:
			// here we start an infinite loop trying to push the record until it is old
			for {
				if time.Since(r.timestamp) > 2*time.Hour {
					log.Printf("dropping old record from %v", r.timestamp)
					break
				} else {
					err := d.insertRecord(r)
					if err != nil {
						log.Printf("failed to insert record at timestamp %v: %s", r.timestamp, err)
						time.Sleep(5 * time.Second)
					} else {
						break
					}
				}
			}
		}

	}
}
