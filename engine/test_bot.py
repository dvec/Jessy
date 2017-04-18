from engine import bot_engine

if __name__ == '__main__':
    while True:
        # WARNING: don't enter commands for bot here!
        print(bot_engine.analyze(None, None, input('> ')))
