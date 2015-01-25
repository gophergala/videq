![videq](http://s3.amazonaws.com/challengepost/photos/production/solution_photos/000/203/847/datas/xlarge.png "videq - High quality video encoding for modern web in golang")

# videq
High quality video encoding for modern web in golang

**Project links and final notes**
   * compo submission page http://gopher-gala.challengepost.com/submissions/31979-videq
   * repo https://github.com/gophergala/videq/
   * demo http://videq.nivas.hr:8080/
     - limitations
       - demo is limited to max 5 min videos and 50 gb of user uploaded data
       - each user can upload only 1 video, after it is encoded or aborted, user can upload new one
       - every hour clenaup will delete old videos

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
Added bonus - it can easily be installed locally for team usage (to save time on uploading to internet server)

## How it works
User drags file to browser (or selects it), and the file is broken into chunks. We start uploading first chunk and right after first chunk is uploaded to the server, server checks if this is a video file suitable for encoding (is it a video file, is it too big, too long etc). If everything is ok we proceeds with upload and encoding.
If upload breaks or stops at some point, user can drag file again and resume. It will be uploaded from last uploaded chunk.
Video encoding is done in queue with help of go workers. First we double pass encode mp4 file which is then base for ogg and webm files. After all files are created, we output it to the user.

**How to use**
*   Download and compile
    go get -u -v https://github.com/gophergala/videq
*   create config file
File has to be present in ./conf/ folder and named eg: ubuntu.config.ini, where ubuntu is the name of server on which it is being run. you get it with uname -n
```
[http]
HOSTNAME = "localhost"
LISTENADDRESS = ":8080"

[db]
HOST="localhost"
NAME="videq"
USER="root"
PASS="XXX"
DEBUG=false
```

*   create and import empty database 
```
mysql -u root -p < sql/create_db.sql
mysql -u root -p videq < sql/empty_db.sql
```
*   run by starting ./start_videq.sh
*   goto http://localhost:8080/


**Prerequisites**
   * server running Ubuntu 14.04 LTS - project was developed on 14.04, possible it could work on anything but we did not have time to test it
   * preinstalled server side applications: 
     * golang (tested with 1.4.1), 
     * mysql (tested on 5.5.40-0ubuntu0.14.04.1 (Ubuntu)), 
     * ffmpeg (tested on 1.2.6-7:1.2.6-1~trusty1), 
     * HandBrakeCLI (0.10.0) 
     * ffmpeg2theora (0.29)

**Tools used in this project**
   * compiler http://golang.org
   * ide      http://www.sublimetext.com/3
   * ffmpeg http://ffmpeg.org/
   * handbrake https://handbrake.fr/
   * ffmpeg2theora http://v2v.cc/~j/ffmpeg2theora/

## Challenges we ran into
- Big file upload
- Agile golang development
- Development of go app on Windows which uses Unix only server side stuff.

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
- move all configurable params to config
- cleanup
- send email when encoding is done
- ...

## The Team
Big shout to fine http://nivas.hr Videq team members:
- https://github.com/guycalledseven backend
- https://github.com/MatejB backend
- https://github.com/luzel frontend
- https://github.com/alencvitkovic design/ux