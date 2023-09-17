"""Integration for Sofar Solar PV inverters with LSW-3 WiFi sticks with SN 23XXXXXXX"""

import platform
import subprocess
import logging
import os

DOMAIN = "sofar_g3_lsw3_logger_reader"
_LOGGER = logging.getLogger(__name__)

def setup(hass, config):
    hass.states.set("sofar_g3_lsw3_logger_reader.status", __name__) # custom_components.sofar_g3_lsw3_logger_reader
    _LOGGER.warning("setup: Working Dir is %s", os.getcwd())        # /config
    _LOGGER.info("setup: Working Dir is %s", os.getcwd())           # main configuration.yaml should enable info, but this is not yet working  

    # When the (go lang) sofar process writes a line to stderr, I'd like it to go to _LOGGER.info(), but I haven't fugured out how to do this
    # without blocking, so for now we connect stdout/err to log files.
    out = open(f"/config/custom_components/{DOMAIN}/out.log", "w") # empty
    err = open(f"/config/custom_components/{DOMAIN}/err.log", "w") # go logging package writes to stderr
    exe = f"/config/custom_components/{DOMAIN}/sofar-arm" if platform.processor.startsWith("arm") else f"/config/custom_components/{DOMAIN}/sofar-x86" 
    proc = subprocess.Popen(exe, cwd=f"/config/custom_components/{DOMAIN}", stdout=out, stderr=err)
    # the (go lang) sofar process reads its config.yaml from cwd, then loops forever
    _LOGGER.warning("setup: end")
    return True
