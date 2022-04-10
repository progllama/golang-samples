import ffmpeg
import argparse

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
 
src = args.src
dst = args.dst
stream = ffmpeg.input(src) 
stream = ffmpeg.output(stream, dst) 
ffmpeg.run(stream)