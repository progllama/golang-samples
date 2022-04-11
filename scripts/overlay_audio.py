# -*- coding: utf-8 -*-
import argparse
import glob
from functools import reduce
from logging import getLogger
from pydub import AudioSegment

parser = argparse.ArgumentParser(description='入力元パターン,出力先を受け取る。')
parser.add_argument(
    "--src",
    required=True
)
parser.add_argument(
    "--dst",
    required=True
)
args = parser.parse_args()

logger = getLogger(__name__)

src = glob.glob(args.src, recursive=True)
dst = args.dst
segments = map(lambda v : AudioSegment.from_file(v, "mp3"), src)
segment = reduce(lambda r, l : l.overlay(r, 0), segments, AudioSegment.empty())
print(segment)
segment.export(dst, format="mp3")