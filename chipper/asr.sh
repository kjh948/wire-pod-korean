#!/bin/bash
rm -rf dump/raw.wav
ffmpeg -f s16le -ar 16k -ac 1 -i dump/raw.pcm dump/raw.wav
python3 asr.py