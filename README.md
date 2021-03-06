# timelapse-go
## Timelapse image temporal linear interpolator

### History

A long time ago I wrote the original timelapse.py, a couple of scripts to allow interpolating an image sequence - for example, in case of dropped frames or simply requiring frames in a timelapse sequence (without dropping the framerate too low).

Original repository: https://github.com/spodzone/timelapse.py

This was a bit over-developed: it had the ability to specify image morphology tweaks at various points in the sequence, relative either to timestamp or source image. For example, one could apply a fade by adjusting the gamma from (1,1,1) to (0,0,0).

Since then, I've reimplemented the core idea (without the morphology) as a timelapse-lite.py script and again from scratch in Julia
Julia repository: https://github.com/spodzone/Timelapse.jl

This project is the third implementation, just the core "lite" functionality, in [Go](https://go.dev/).

### Building

```
$ git clone https://github.com/spodzone/timelapse-go.git
$ go build
```

### Running

Assuming you have a lot of JPEG images in `indir` and an empty directory `outdir` ready to receive `noframes` number of interpolated frames,

`$ ./timelapse.exe indir noframes outdir`

`$ ffmpeg -i outdir/img%05d.jpg -y -qscale 0 -r 25 timelapse.flv `

That should be it.
