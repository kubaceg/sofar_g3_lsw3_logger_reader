package sofar

type field struct {
	register  int
	name      string
	valueType string
	factor    float32
	unit      string
}

type registerRange struct {
	start       int
	end         int
	replyFields []field
}

var allRegisterRanges = []registerRange{
	rrGridOutput,
	rrPVOutput,
	rrPVTotal,
	rrEnergyTodayTotals,
	rrSystemInfo,
	rrSystemHardware,
	rrBatOutput,
	rrRatio,
}

func GetAllRegisterNames() []string {
	result := make([]string, 0)
	for _, rr := range allRegisterRanges {
		for _, f := range rr.replyFields {
			if f.name == "" || f.valueType == "" {
				// Measurements without a name or value type are ignored in replies
				continue
			}
			result = append(result, f.name)
		}
	}
	return result
}

var rrSystemInfo = registerRange{
	start: 0x400,
	end:   0x43A,
	replyFields: []field{
		{0x0404, "Status/SysState", "U16", 0, ""},
		{0x0405, "Status/Fault01", "U16", 0, ""},
		{0x0406, "Status/Fault02", "U16", 0, ""},
		{0x0407, "Status/Fault03", "U16", 0, ""},
		{0x0408, "Status/Fault04", "U16", 0, ""},
		{0x0409, "Status/Fault05", "U16", 0, ""},
		{0x040A, "Status/Fault06", "U16", 0, ""},
		{0x040B, "Status/Fault07", "U16", 0, ""},
		{0x040C, "Status/Fault08", "U16", 0, ""},
		{0x040D, "Status/Fault09", "U16", 0, ""},
		{0x040E, "Status/Fault10", "U16", 0, ""},
		{0x040F, "Status/Fault11", "U16", 0, ""},
		{0x0410, "Status/Fault12", "U16", 0, ""},
		{0x0411, "Status/Fault13", "U16", 0, ""},
		{0x0412, "Status/Fault14", "U16", 0, ""},
		{0x0413, "Status/Fault15", "U16", 0, ""},
		{0x0414, "Status/Fault16", "U16", 0, ""},
		{0x0415, "Status/Fault17", "U16", 0, ""},
		{0x0416, "Status/Fault18", "U16", 0, ""},
		{0x0417, "Status/Countdown", "U16", 1, "seconds"},
		{0x0418, "Hardware/Temperature/Env1", "I16", 1, "C"},
		{0x0419, "Hardware/Temperature/Env2", "I16", 1, "C"},
		{0x041A, "Hardware/Temperature/HeatSink1", "I16", 1, "C"},
		{0x041B, "Hardware/Temperature/HeatSink2", "I16", 1, "C"},
		{0x041C, "Hardware/Temperature/HeatSink3", "I16", 1, "C"},
		{0x041D, "Hardware/Temperature/HeatSink4", "I16", 1, "C"},
		{0x041E, "Hardware/Temperature/HeatSink5", "I16", 1, "C"},
		{0x041F, "Hardware/Temperature/HeatSink6", "I16", 1, "C"},
		{0x0420, "Hardware/Temperature/Inv1", "I16", 1, "C"},
		{0x0421, "Hardware/Temperature/Inv2", "I16", 1, "C"},
		{0x0422, "Hardware/Temperature/Inv3", "I16", 1, "C"},
		{0x0423, "Hardware/Temperature/Rsvd1", "I16", 1, "C"},
		{0x0424, "Hardware/Temperature/Rsvd2", "I16", 1, "C"},
		{0x0425, "Hardware/Temperature/Rsvd3", "I16", 1, "C"},
		{0x0426, "Time/Generation/Today", "U16", 1, "Minute"},
		{0x0427, "Time/Generation/Total", "U32", 1, "Minute"},
		{0x0429, "Time/Service/Total", "U32", 1, "Minute"},
		{0x042B, "Status/InsulationResistance", "U16", 1, "kOhm"},
		{0x0432, "Status/Fault19", "U16", 0, ""},
		{0x0433, "Status/Fault20", "U16", 0, ""},
		{0x0434, "Status/Fault21", "U16", 0, ""},
		{0x0435, "Status/Fault22", "U16", 0, ""},
		{0x0436, "Status/Fault23", "U16", 0, ""},
		{0x0437, "Status/Fault24", "U16", 0, ""},
		{0x0438, "Status/Fault25", "U16", 0, ""},
		{0x0439, "Status/Fault26", "U16", 0, ""},
		{0x043A, "Status/Fault27", "U16", 0, ""},
	},
}

var rrSystemHardware = registerRange{
	start: 0x440,
	end:   0x44E,
	replyFields: []field{
		{0x0444, "Hardware/Production_Code", "U16", 0, ""},
		{0x044D, "Hardware/Hardware_Version0", "", 0, ""},
		{0x044E, "Hardware/Hardware_Version1", "", 0, ""},
	},
}

var rrEnergyTodayTotals = registerRange{
	start: 0x680,
	end:   0x69B,
	replyFields: []field{
		{0x0684, "Energy/PV_Generation_Today", "U32", 0.01, "kWh"},
		{0x0686, "Energy/PV_Generation_Total", "U32", 0.01, "kWh"},
		{0x0688, "Energy/Load_Consumption_Today", "U32", 0.01, "kWh"},
		{0x068A, "Energy/Load_Consumption_Total", "U32", 0.1, "kWh"},
		{0x068C, "Energy/Energy_Purchase_Today", "U32", 0.01, "kWh"},
		{0x068E, "Energy/Energy_Purchase_Total", "U32", 0.1, "kWh"},
		{0x0690, "Energy/Energy_Selling_Today", "U32", 0.01, "kWh"},
		{0x0692, "Energy/Energy_Selling_Total", "U32", 0.1, "kWh"},
		{0x0694, "Energy/Bat_Charge_Today", "U32", 0.01, "kWh"},
		{0x0696, "Energy/Bat_Charge_Total", "U32", 0.1, "kWh"},
		{0x0698, "Energy/Bat_Discharge_Today", "U32", 0.01, "kWh"},
		{0x069A, "Energy/Bat_Discharge_Total", "U32", 0.1, "kWh"},
	},
}

var rrPVOutput = registerRange{
	start: 0x0580,
	end:   0x0589,
	replyFields: []field{
		{0x0584, "PV/String1/Voltage", "U16", 0.1, "V"},
		{0x0585, "PV/String1/Current", "U16", 0.1, "A"},
		{0x0586, "PV/String1/Power", "U16", 0.01, "kW"},
		{0x0587, "PV/String2/Voltage", "U16", 0.1, "V"},
		{0x0588, "PV/String2/Current", "U16", 0.1, "A"},
		{0x0589, "PV/String2/Power", "U16", 0.01, "kW"},
	},
}

var rrPVTotal = registerRange{
	start: 0x05C0,
	end:   0x05C4,
	replyFields: []field{
		{0x05C3, "PV/Total/Power", "U16", 0.1, "kW"},
	},
}

var rrGridOutput = registerRange{
	start: 0x480,
	end:   0x4bc,
	replyFields: []field{
		{0x0484, "Grid/Output/Frequency", "U16", 0.01, "Hz"},
		{0x0485, "Grid/Output/Total_ActivePower", "I16", 0.01, "kW"},
		{0x0486, "Grid/Output/Total_ReactivePower", "I16", 0.01, "kW"},
		{0x0487, "Grid/Output/Total_ApparentPower", "I16", 0.01, "kW"},
		{0x0488, "Grid/PCC/ActivePower", "I16", 0.01, "kW"},
		{0x0489, "Grid/PCC/ReactivePower", "I16", 0.01, "kW"},
		{0x048A, "Grid/PCC/ApparentPower", "I16", 0.01, "kW"},
		{0x048B, "Grid/Output/Rsvd1_Grid", "", 0, ""},
		{0x048C, "Grid/Output/Rsvd2_Grid", "", 0, ""},
		{0x048D, "Grid/Output/R_Voltage", "U16", 0.1, "V"},
		{0x048E, "Grid/Output/R_Current", "U16", 0.01, ""},
		{0x048F, "Grid/Output/R_ActivePower", "I16", 0.01, "kW"},
		{0x0490, "Grid/Output/R_ReactivePower", "I16", 0.01, "kW"},
		{0x0491, "Grid/Output/R_PowerFactor", "I16", 0.001, "p.u."},
		{0x0492, "Grid/PCC/R_Current", "U16", 0.01, ""},
		{0x0493, "Grid/PCC/R_ActivePower", "I16", 0.01, "kW"},
		{0x0494, "Grid/PCC/R_ReactivePower", "I16", 0.01, "kW"},
		{0x0495, "Grid/PCC/R_PowerFactor", "I16", 0.001, "p.u."},
		{0x0496, "Grid/Output/R_Rsvd1", "", 0, ""},
		{0x0497, "Grid/Output/R_Rsvd2", "", 0, ""},
		{0x0498, "Grid/Output/S_Voltage", "U16", 0.1, "V"},
		{0x0499, "Grid/Output/S_Current", "U16", 0.01, ""},
		{0x049A, "Grid/Output/S_ActivePower", "I16", 0.01, "kW"},
		{0x049B, "Grid/Output/S_ReactivePower", "I16", 0.01, "kW"},
		{0x049C, "Grid/Output/S_PowerFactor", "I16", 0.001, "p.u."},
		{0x049D, "Grid/PCC/S_Current", "U16", 0.01, ""},
		{0x049E, "Grid/PCC/S_ActivePower", "I16", 0.01, "kW"},
		{0x049F, "Grid/PCC/S_ReactivePower", "I16", 0.01, "kW"},
		{0x04A0, "Grid/PCC/S_PowerFactor", "I16", 0.001, "p.u."},
		{0x04A1, "Grid/Output/S_Rsvd1", "", 0, ""},
		{0x04A2, "Grid/Output/S_Rsvd2", "", 0, ""},
		{0x04A3, "Grid/Output/T_Voltage_Phase_T", "U16", 0.1, "V"},
		{0x04A4, "Grid/Output/T_Current", "U16", 0.01, ""},
		{0x04A5, "Grid/Output/T_ActivePower", "I16", 0.01, "kW"},
		{0x04A6, "Grid/Output/T_ReactivePower", "I16", 0.01, "kW"},
		{0x04A7, "Grid/Output/T_PowerFactor", "I16", 0.001, "p.u."},
		{0x04A8, "Grid/PCC/T_Current", "U16", 0.01, ""},
		{0x04A9, "Grid/PCC/T_ActivePower", "I16", 0.01, "kW"},
		{0x04AA, "Grid/PCC/T_ReactivePower", "I16", 0.01, "kW"},
		{0x04AB, "Grid/PCC/T_PowerFactor", "I16", 0.001, "p.u."},
		{0x04AC, "Grid/Output/T_Rsvd1", "", 0, ""},
		{0x04AD, "Grid/Output/T_Rsvd2", "", 0, ""},
		{0x04AE, "Grid/Output/ActivePower_PV_Ext", "U16", 0.01, "kW"},
		{0x04AF, "Grid/Output/ActivePower_Load_Sys", "U16", 0.01, "kW"},
		{0x04B0, "Grid/Output/L1N_Voltage_Phase_L1N", "U16", 0.1, "V"},
		{0x04B1, "Grid/Output/L1N_Current", "U16", 0.01, ""},
		{0x04B2, "Grid/Output/L1N_ActivePower", "I16", 0.01, "kW"},
		{0x04B3, "Grid/PCC/L1N_Current", "U16", 0.01, ""},
		{0x04B4, "Grid/PCC/L1N_ActivePower", "I16", 0.01, "kW"},
		{0x04B5, "Grid/Output/L2N_Voltage_Phase_L2N", "U16", 0.1, "V"},
		{0x04B6, "Grid/Output/L2N_Current", "U16", 0.01, ""},
		{0x04B7, "Grid/Output/L2N_ActivePower", "I16", 0.01, "kW"},
		{0x04B8, "Grid/PCC/L2N_Current", "U16", 0.01, ""},
		{0x04B9, "Grid/PCC/L2N_ActivePower", "I16", 0.01, "kW"},
		{0x04BA, "Grid/Output/Voltage_Line_L1", "U16", 0.1, "V"},
		{0x04BB, "Grid/Output/Voltage_Line_L2", "U16", 0.1, "V"},
		{0x04BC, "Grid/Output/Voltage_Line_L3", "U16", 0.1, "V"},
	},
}

var rrBatOutput = registerRange{
	start: 0x600,
	end:   0x611,
	replyFields: []field{
		{0x0604, "Battery/1/Voltage_Bat", "U16", 0.1, "V"},
		{0x0605, "Battery/1/Current_Bat", "I16", 0.01, "A"},
		{0x0606, "Battery/1/Power_Bat", "I16", 0.01, "kW"},
		{0x0607, "Battery/1/Temperature_Env_Bat", "I16", 1, "C"},
		{0x0608, "Battery/1/SOC_Bat", "U16", 1, "%"},
		{0x0609, "Battery/1/SOH_Bat", "U16", 1, "%"},
		{0x060A, "Battery/1/ChargeCycle_Bat", "U16", 1, ""},
		{0x060B, "Battery/2/Voltage_Bat", "U16", 0.1, "V"},
		{0x060C, "Battery/2/Current_Bat", "I16", 0.01, "A"},
		{0x060D, "Battery/2/Power_Bat", "I16", 0.01, "kW"},
		{0x060E, "Battery/2/Temperature_Env_Bat", "I16", 1, "C"},
		{0x060F, "Battery/2/SOC_Bat", "U16", 1, "%"},
		{0x0610, "Battery/2/SOH_Bat", "U16", 1, "%"},
		{0x0611, "Battery/2/ChargeCycle_Bat", "U16", 1, ""},
	},
}

var rrRatio = registerRange{
	start: 0x1030,
	end:   0x103D,
	replyFields: []field{
		{0x1039, "Ratio/PV_Generation_Ratio", "U16", 0.001, ""},
		{0x103A, "Ratio/Energy_Purchase_Ratio", "U16", 0.001, ""},
		{0x103B, "Ratio/Energy_Selling_Ratio", "U16", 0.001, ""},
		{0x103C, "Ratio/Bat_Charge_Ratio", "U16", 0.001, ""},
		{0x103D, "Ratio/Bat_Discharge_Ratio", "U16", 0.001, ""},
	},
}
