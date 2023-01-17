# webpfex
Extract frames from an animated WEBP or convert them to MP4. Relies on webpmux and ffmpeg.

## Why?
Tools like ImageMagick can extract animated WEBP frames, but said frames are extracted directly as-is as stored in the WEBP file. Some animated WEBP files only store successive changes from previous frames, thus have transparency or are of different resolution. These frames can't just be extracted and fed into programs like FFmpeg to reconstruct them as a video or other animated image formats. This program fixes that.