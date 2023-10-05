# Grafana dashboard

Grafana dashboard for sofar inverter data exported [sofar_g3_lsw3_logger_reader](https://github.com/kubaceg/sofar_g3_lsw3_logger_reader). Just paste [this json config](grafana_dashboard.json) in Your grafana and make some modifications depending on your config:

* replace all `8000` int occurences in json to maximum power of Your insalation in watt.
* this config uses `sofar_logger_` metric prefix, if You use another replace it in every query.

![Alt text](relative%20dashboard.png?raw=true "Grafana dashboard")
