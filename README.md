# gopush

Simple Pushbullet API client and CLI tool

This is not really feature complete, but works well enough for me. Use at your own risk. I thought I was going to build it out more, but basically never got around to it due to lack of needing new features.

General usage:
```
gopush -token="yourapitokenfrompushbullet"
gopush "Whatever you want pushed"
```

Defaults to all devices, but you can also specify with `gopush -d="Device name"`. Can't remember your devices? `gopush -l` will list them.

I often use this combined with a bash or fish script (maybe I'll post later) to push me successes, failures, or results for long running commands.
