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

var UploadLogic = {
	firstPart: false,
	onProgress: function() {
		if flow.progress() == 0 {
			return
		}

		if UploadLogic.firstPart == false {
			flow.pause();
			$.ajax({
				url: "/check/",
				dataType: "json"
			}).done(function(data) {
				console.log(data)
			});
		}
		else {
			var progress = flow.progress() * 100;
			if (progress < 100) {
				$('.progress-bar').css('width', progress + '%');
			}
		}
	}
}


$(function(){

	var notCompletedFiles = getFilesListFromCookie();
	if (notCompletedFiles != false) {
		$.each(notCompletedFiles, function(ix, fileName){
			$('#list-of-not-uploded-files').append('<li>' + fileName + '</li>');
		});
		$('#list-of-not-uploded-files-holder').show();
	}

	var flow = new Flow({
	  target:'/upload/',
	  uploadMethod: 'POST',
	  testChunks: true,
	  simultaneousUploads: 1,
	  prioritizeFirstAndLastChunk: true,
	  progressCallbacksInterval: 0
	});
	if (!flow.support) {
		alert("Browser dose not support modern upload!");
		$('#upload_form').hide();
	}

	flow.assignBrowse(document.getElementById('js-upload-files'));
	flow.assignDrop(document.getElementById('drop-zone'));

	flow.on('fileAdded', function(file, event){
		addToCookie(file.name);
	    $('#fileLog').append('<a href="#" class="list-group-item list-group-item">' + file.name + ' added to upload queue</a>');
	});
	flow.on('fileSuccess', function(file,message){
		removeFromCookie(file.name);
	    $('#fileLog').append('<a href="#" class="list-group-item list-group-item-success"><span class="badge alert-success pull-right">Success</span>' + file.name + ' ' + message + '</a>');
	});
	flow.on('fileError', function(file, message){
		removeFromCookie(file.name);
	    $('#fileLog').append('<a href="#" class="list-group-item list-group-item-danger"><span class="badge alert-danger pull-right">Error</span>' + file.name + ' ' + message + '</a>');
	});

	flow.on('progress', UploadLogic.onProgress);

	flow.on('complete', function(){
	    $('.progress-bar').css('width', '0%');
	    $('#fileLog').append('<a href="#" class="list-group-item list-group-item-success"><span class="badge alert-success pull-right">Success</span>All upload completed</a>');
	    $('#upload_form').show();
	});
	flow.on('uploadStart', function(){
	    $('#upload_form').hide();
		$('#list-of-not-uploded-files-holder').hide();
	});

	$('#js-upload-form').submit(function(ev){
		ev.preventDefault();
		flow.upload();
	});

});