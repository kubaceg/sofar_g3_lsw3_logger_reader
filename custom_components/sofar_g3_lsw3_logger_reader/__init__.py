"""Integration for Sofar Solar PV inverters with LSW-3 WiFi sticks with SN 23XXXXXXX"""

import platform
import subprocess
import logging
import os

_LOGGER = logging.getLogger(__name__)
_DOMAIN = "sofar_g3_lsw3_logger_reader"
_DIR = f"/config/custom_components/{_DOMAIN}"

def setup(hass, config):
    hass.states.set(f"{_DOMAIN}.status", __name__)                  # custom_components.sofar_g3_lsw3_logger_reader
    _LOGGER.warning("setup: Working Dir is %s", os.getcwd())        # /config
    _LOGGER.info("setup: Working Dir is %s", os.getcwd())           # main configuration.yaml should enable info, but this is not yet working  

    # When the process writes a line to stderr, I'd like it to go to _LOGGER.info(), but I haven't fugured out how to do this
    # without blocking, so for now we connect stdout/err to log files.
    out = open(f"{_DIR}/out.log", "w")                              # empty
    err = open(f"{_DIR}/err.log", "w")                              # go logging package writes to stderr
    # platform.processor() is '' on HA OS on Raspberry Pi
    exe = f"{_DIR}/sofar-x86" if platform.processor().startswith("x86") else f"{_DIR}/sofar-arm" 
    proc = subprocess.Popen(exe, cwd=f"{_DIR}", stdout=out, stderr=err)
    # the (go lang) sofar process reads its config.yaml from cwd, then loops forever
    _LOGGER.warning("setup: end")
    return True
