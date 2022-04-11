import ffmpeg
import argparse
import glob
from logging import getLogger

parser = argparse.ArgumentParser(description='入力元,出力先を受け取る。')
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

try: 
    src = args.src
    dst = args.dst
    stream = ffmpeg.input(src) 
    stream = ffmpeg.output(stream, dst) 
    ffmpeg.run(stream)
except:
    logger.info("エラーを起こしました。考えられる原因を以下に記載します。")
    logger.info("・そもそも動画に音声がない。")