from unittest import TestCase
from engine import functions


class TestFunctions(TestCase):
    def test_get_random_num(self):
        tests = (
            ('test1', 'test2'),
            ('qwerty', 'qwerty.'),
            ('12345', '12346'),
            ('!@#$%^&*()_+|\\/"][}{', '!&(*@#&!)*')
        )
        for test1, test2 in tests:
            self.assertNotEquals(functions.get_random_num(test1), functions.get_random_num(test2))

    def test_to_simple_text(self):
        tests = (
            ('', '!?)(09:.'),
            ('qwerty', 'qwerty.'),
            ('HELLO', 'hello'),
            ('JuSt A tEsT', 'just a test')
        )
        for test1, test2 in tests:
            self.assertEquals(functions.to_simple_text(test1), functions.to_simple_text(test2))
