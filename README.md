# videq
High quality video encoding for modern web in golang

## Inspiration

We work in digital agency and we build a lot of websites that use videos (full screen, background, interviews, product showcase etc.). 
Developing compatible cross platform/browser website is a technical nightmare due to thrilling mix of “historical” and "future"  standards. In relation to video playback - that means if we want a video on website to be playable on most platforms we need ad least 3 different video formats:
- H.264+AAC+MP4 - Safari, Chrome, Firefox (new versions), IE
- Theora+Vorbis+Ogg – Firefox
- WebM - Chrome, IE
 
Each format has its own limitations and settings that must be applied in order to get most compatible and excellent looking output, using free converters available to general public.

## The problem(s)
In order to convert one video to specified output formats you need different multidisciplinary knowledge, different desktop software, time and patience. When you are building something under deadline, there is no time and patience, so we wanted it to make it simple for anybody to use (but it’s mostly for web developers).

Uploads of big files via web browser have always been one of big issues on the web.  Hopefully we managed to solve ti by breaking the file in chunks and handling upload of each chunk.
 
## Solution
An UX friendly, easy to use website/webapp for re-coding video files to „safe for web“ video formats. User drags videos, waits, downloads converted videos+html+fallback static image/anim gif.

## How it works
User drags file to browser.
After file is broken into chunks, we start uploading first chunk. Right after first chunk is uploaded to the server, server checks if this is a video file suitable for encoding. If it is - it proceeds with upload and encoding.
If upload breaks/stops at some point, user can drag file again and resume. It will be uploaded from last uploaded chunk.
Encoding part is done in queue with help of go workers. First we double pass encode mp4 file which is then base for ogg and webm files. After all files are created, we output it to the user.

## Challenges we ran into
Big file upload
Development of go app on Windows which uses Unix only server side stuff.

## Accomplishments that we are proud of
We managed to get MVP in just two days. :) 

## What we learned
We are new in golang and  we used this hackaton to learn more about: 
- chunked uploads handling (for extreme large files) 
- team workflow on short deadline golang projects


## What's next for videq
- more encoding options for quality
- graceful restart
- scalable (multi server workers) … if we manage in that short time we have… 

## The Team
Big shout to fine http://nivas.hr Videq team members:
- https://github.com/guycalledseven backend
- https://github.com/MatejB backend
- https://github.com/luzel frontend
- https://github.com/alencvitkovic design/ux
