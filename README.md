# Paprika

Paprika is a toolbox for creating short clips from multiple png files.

<p>
<img src="https://raw.githubusercontent.com/alvoras/paprika/master/media/clip_export.gif" />
</p>

# Requirements

This program is mainly a wrapper around ffmpeg. As such, you must have ffmpeg in your PATH.

> :warning: The crossfade effect only works with ffmpeg >= 4.4

# Installation

Either build the binary from source or grab the latest release.

If you have a working Golang installation, you can build from source :

```
go get
go build .
```

You can also install it with the following command if your Golang's bin directory is in your PATH :

```
go install .
```

# Quickstart

The program has multiple subcommands.

All the examples below will use this as a source :

```
$ ls -l ./frames                                                                                                              
0000.png
0001.png
0002.png
0003.png
...
0449.png
```

> :warning: The frames files **MUST** be in the form **NNNN.png, eg. "0000.png, 0001.png, [...], 9999.png"**
> 
> They also must be continuous, as the missing numbers will result in black frames in the final clip
>
> Keep in mind that the first frame is 0000.png

## Clip

The "clip" subcommand is used to assemble multiple png files into a single mp4 video.

Here are the options :

```
Usage:
  paprika clip <path> [flags]

Flags:
  -b, --boomerang            Apply a boomerang effect
  -c, --cut string           Shorthand for the --start-at and --end-at combo. Syntax : '--cut start:end'. Example to make a clip from the 10th to the 20th frames : '--cut 10:20'
  -S, --start-at int         Start at the specified frame
  -E, --end-at int           End at the specified frame
  -f, --first int            Use the first N frames
  -l, --last int             Use the last N frames
  -F, --fps int              Frame per second (default 24)
  -g, --gif                  Convert to gif instead of mp4
  -d, --gif-delay int        Delay to apply for each frame (gif only) (default 4)
  -h, --help                 help for clip
  -o, --out string           Path to save the clip to
  -x, --xfade                Apply a crossfade effect
  -D, --xfade-duration int   Specify the duration in milliseconds of the crossfade effect. Default is half the total duration
  -O, --xfade-offset int     Specify the time offset in milliseconds at which the crossfade starts. Default is a quarter of the total duration
```

In its most simple form, this command will assemble all the frames present in `./frames` at 24 fps (by default) and output a `frames.mp4` video file :

```
paprika clip ./frames
```

### Editing
You can specify where to cut from the sources using these flags :

```
  -c, --cut string           Shorthand for the --start-at and --end-at combo. Syntax : '--cut start:end'. Example to make a clip from the 10th to the 20th frames : '--cut 10:20'
  -S, --start-at int         Start at the specified frame
  -E, --end-at int           End at the specified frame
  -f, --first int            Use the first N frames
  -l, --last int             Use the last N frames
```

This will export a clip from the frame 75 to 150 at 25 fps :

```
paprika clip ./frames --cut 75:150 --fps 25
```

This will export a clip from the frame 75 to the end :

```
paprika clip ./frames --first 75
```

This will export the last 50 frames :

```
paprika clip ./frames --last 50
```

This will export a clip starting at the 75th frame to the end :

```
paprika clip ./frames --start-at 75
```

This will export a clip from the beginning to the 150th frame :

```
paprika clip ./frames --end-at 150
```

This will export a clip starting at the 75th frame to the 150th (equivalent to `--cut 75:150`) :

```
paprika clip ./frames --start-at 75 --end-at 150
```

### Effects

A few effects are available.

#### Crossfade

Crossfade is enabled with the `-x|--xfade` flag.

The crossfade effect is useful to make pseudo-seamlessly looping videos. It will split the clip in half, put the second part in front of the first and apply a crossfade between the two.

By default, the crossfade is applied starting at a quarter of the video length, for a duration of half its length minus 1sec.

You can override the default effect's parameters with `--xfade-offset|O` and `--xfade-duration|D`. The values must be in milliseconds.

Example : 

With default parameters :

```
paprika clip ./frames --xfade 
```

To put a 1 second crossfade after 2 seconds of video :

```
paprika clip ./frames --cut 150:250 --xfade --xfade-offset 2000 --xfade-duration 1000
```

#### Boomerang

Classic boomerang effect. A reversed copy of the video will be appended at the end of the clip.

This effect is enabled with the `-b|--boomerang` flag.

Example : 

```
paprika clip ./frames --boomerang
```

### Gif

It is also possible to export the clip as a gif. Keep in mind that this gif will not be compressed and thus very heavy. It can also take a lot of memory while building it.

> :warning: The crossfade effect is incompatible with the gif mode 

## Extract

An extract subcommand is available as well.
<p>
<img src="https://raw.githubusercontent.com/alvoras/paprika/master/media/extract.png" />
</p>

```
Usage:
  paprika extract <path> [flags]

Flags:
  -h, --help         help for extract
  -o, --out string   Output directory (default "./extracted")
  -s, --step int     Save an image every N pictures (default 50)
```

This is useful when curating stepped generative work, such as AI generative art based on iterations (VQGAN+CLIP, CLIP guided diffusion...).

By default, this will extract a picture every 50 at `./extracted/frames/` :

```
paprika extract ./frames
```

This will do the same, except with a step of 75 frames :

```
paprika extract ./frames -s 75
```
