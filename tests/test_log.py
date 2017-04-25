from unittest import TestCase
from log.log import Log


class TestLog(TestCase):
    def test_init_log(self):
        Log.init_log()
