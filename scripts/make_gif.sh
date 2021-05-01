#!/bin/sh
# Temp files
mp4="recording/temp/out.mp4"
palette="recording/temp/palette.png"

gif="recording/output/$1.gif"

fps="50"

# Create mp4 from png's
ffmpeg -r $fps -start_number 0 -i recording/temp/image%03d.png -c:v libx264 -vf "fps=$fps,format=yuv420p" $mp4

# Create palette, so that we can have high quality gifs at low cost
ffmpeg -v warning -i $mp4 -vf "fps=$fps,palettegen" -y $palette

# Convert mp4 to gif
ffmpeg -v warning -i $mp4 -i $palette -lavfi "fps=$fps [x]; [x][1:v] paletteuse" -y $gif

# Remove working files
rm -rf recording/temp