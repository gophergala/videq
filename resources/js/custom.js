function addToCookie(filename) {
	removeFromCookie(filename)

	var cookieStorage = $.cookie('files');
	if (cookieStorage == undefined) {
		cookieStorage = '';
	}

	cookieStorage = cookieStorage + filename + '|';

	$.cookie('files', cookieStorage, { expires: 7, path: '/' });
}

function removeFromCookie(filename) {
	var cookieStorage = $.cookie('files');
	if (cookieStorage == undefined) {
		cookieStorage = '';
		$.cookie('files', cookieStorage, { expires: 7, path: '/' });
		return;
	}

	cookieStorage = cookieStorage.replace(filename + '|',"");

	$.cookie('files', cookieStorage, { expires: 7, path: '/' });
}

function getFilesListFromCookie() {
	var cookieStorage = $.cookie('files');
	if (cookieStorage == undefined || cookieStorage == ""){
		return false;
	}

	var parts = cookieStorage.split("|");
	parts.pop();

	return parts
}

function refreshListOfStoredFiles() {
	$.ajax({
		url: "/list-of-files/",
		dataType: "json"
	}).done(function(files) {
		$.each(files, function(ix, fileName){
			$('#list-of-stored-files').append('<li><a href="/download/?filename=' + encodeURIComponent(fileName) + '">' + fileName + '</a></li>');
		});

		if (files.length > 0) {
			$('#list-of-stored-files-holder').show();
		} else {
			$('#list-of-stored-files-holder').hide();
		}
	});
}


var Videq = function () {


	this.init = function () {

		/* ------------------------------------------------------ */
		/* Homepage / Add video to intro block
		/* ------------------------------------------------------ */
		var root = $('body.home'),
			vid_cont = $('#intro_vid', root);

		if( vid_cont.length > 0 && !Modernizr.touch && Modernizr.mq('only screen and (min-width: 781px)') ) {


			var html = '<video autoplay loop poster="/resources/vid/loop.jpg">'+
					   '<source src="/resources/vid/loop.mp4" type="video/mp4" />'+
					   '<source src="/resources/vid/loop.webm" type="video/webm" />'+
					   '<source src="/resources/vid/loop.ogv" type="video/ogg" />'+
					   '</video>';

			vid_cont.prepend(html);

			var video = $('video', vid_cont).get(0);

			$(window).scroll(function() {
				pauseVid();
			});

			function pauseVid() {
				if( $(window).scrollTop() > vid_cont.height() ) {
					video.pause();
				} else {
					if( video.paused ) {
						video.play();
					}
				}
			}

			// init
			pauseVid();

		}

	};


	this.submit = function () {

	};

	return {
		init: this.init
	};

}



$(function(){


	var ggg = Videq();
	ggg.init();


	refreshListOfStoredFiles();

	var notCompletedFiles = getFilesListFromCookie();
	if (notCompletedFiles != false) {
		$.each(notCompletedFiles, function(ix, fileName){
			$('#list-of-not-uploded-files').append('<li>' + fileName + '</li>');
		});
		$('#list-of-not-uploded-files-holder').show();
	}

	var flow = new Flow({
	  target: '/upload/',
	  uploadMethod: 'POST',
	  testChunks: true,
	  singleFile: true
	});

	if (!flow.support) {
		alert("Browser dose not support modern upload!");
		$('#upload_form').hide();
	}

	flow.assignBrowse(document.getElementById('js-upload-files'));
	flow.assignDrop(document.getElementById('drop-zone'));

	flow.on('fileAdded', function(file, event){
		addToCookie(file.name);
		flow.upload();
	});

	flow.on('fileSuccess', function(file,message){
		removeFromCookie(file.name);
	    $('#fileLog').append('<a href="#" class="list-group-item list-group-item-success"><span class="badge alert-success pull-right">Success</span>' + file.name + ' ' + message + '</a>');
	});
	flow.on('fileError', function(file, message){
		removeFromCookie(file.name);
	    $('#fileLog').append('<a href="#" class="list-group-item list-group-item-danger"><span class="badge alert-danger pull-right">Error</span>' + file.name + ' ' + message + '</a>');
	});
	flow.on('progress', function(){
		var progress = flow.progress() * 100;
		if (progress < 100) {
			$('.progress-bar').css('width', progress + '%');
		}
	});
	flow.on('complete', function(){
	    $('.progress-bar').css('width', '0%');
	    $('#fileLog').append('<a href="#" class="list-group-item list-group-item-success"><span class="badge alert-success pull-right">Success</span>All upload completed</a>');
	    $('#upload_form').show();
	    refreshListOfStoredFiles();
	});
	flow.on('uploadStart', function(){
	    $('#upload_form').hide();
		$('#list-of-not-uploded-files-holder').hide();
	});

	$('#js-upload-form').submit(function(ev){
		ev.preventDefault();
		flow.upload();
	});

	$('.classic_upload').on("click", function(ev){
		$('input[type=file]').click();
    	return false;
	});

	$(".bswitch").bootstrapSwitch();

});