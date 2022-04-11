import ffmpeg
import argparse
import os
import glob
from logging import getLogger

parser = argparse.ArgumentParser(description='入力元パターン(mp4),出力先ディレクトリを受け取る。')
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
print("---------------------処理開始---------------------")

try: 
    src = glob.glob(args.src, recursive=True)
    src_and_dst = map(lambda v : (v, os.path.join(args.dst, os.path.basename(v)[:-1]+"3")), src)
    in_streams_and_dst = [(ffmpeg.input(s), d) for s, d in src_and_dst]
    out_streams = [ffmpeg.output(ist, d) for ist, d in in_streams_and_dst]
    for stream in out_streams:
        ffmpeg.run(stream)
except:
    logger.info("エラーを起こしました。考えられる原因を以下に記載します。")
    logger.info("・そもそも動画に音声がない。")