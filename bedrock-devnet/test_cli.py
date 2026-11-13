import importlib
import unittest


class DevnetCLITest(unittest.TestCase):
    def test_boolean_flags(self):
        try:
            parser = importlib.import_module('devnet').parser
        except TypeError as err:
            self.fail(f'Devnet CLI must initialize on supported Python versions: {err}')

        for flag in ('allocs', 'test', 'init'):
            with self.subTest(flag=flag):
                self.assertIsNone(getattr(parser.parse_args([]), flag))
                self.assertIs(getattr(parser.parse_args([f'--{flag}']), flag), True)
                self.assertIs(getattr(parser.parse_args([f'--no-{flag}']), flag), False)


if __name__ == '__main__':
    unittest.main()
