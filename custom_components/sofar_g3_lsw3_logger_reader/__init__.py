"""Integration for Sofar Solar PV inverters with LSW-3 WiFi sticks with SN 23XXXXXXX"""

import subprocess
import logging
import os

DOMAIN = "sofar_g3_lsw3_logger_reader"
_LOGGER = logging.getLogger(__name__)

def setup(hass, config):
    hass.states.set("sofar_g3_lsw3_logger_reader.status", __name__)
    _LOGGER.warning("setup: Working Dir is %s", os.getcwd())
    _LOGGER.info("setup: Working Dir is %s", os.getcwd())  # main configuration.yaml should enable info, but not yet working  

    out = open(f"/config/custom_components/{DOMAIN}/out.log", "w") # empty
    err = open(f"/config/custom_components/{DOMAIN}/err.log", "w") # go logging package writes to stderr
    proc = subprocess.Popen(f"/config/custom_components/{DOMAIN}/sofar", cwd=f"/config/custom_components/{DOMAIN}", stdout=out, stderr=err)
    # the sofar process is left running forever
    _LOGGER.warning("setup: end")
    return True
