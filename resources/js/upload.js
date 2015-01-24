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

/**
 * Flow JS
 */
var UploadLogic = {

	flow : false,
	firstPart: false,
	
	init : function () {

		UploadLogic.flow = new Flow({
			target: '/upload/',
			uploadMethod: 'POST',
			testChunks: true,
			simultaneousUploads: 1,
			prioritizeFirstAndLastChunk: true,
			progressCallbacksInterval: 0,
			singleFile: true
		});


		if (!UploadLogic.flow.support) {
			alert("Browser dose not support modern upload!");
			$('#upload_form').hide();
		}

		UploadLogic.bindEvents();
	},

	bindEvents: function () {

		var self = this;

		UploadLogic.flow.assignBrowse(document.getElementById('js-upload-files'));
		UploadLogic.flow.assignDrop(document.getElementById('drop-zone'));

		UploadLogic.flow.on('fileAdded', function(file, event){
			addToCookie(file.name);

			setTimeout(function(){

				UploadLogic.flow.upload();
				
			},1);

		});

		UploadLogic.flow.on('fileSuccess', function(file,message){
			removeFromCookie(file.name);
		    $('#fileLog').append('<a href="#" class="list-group-item list-group-item-success"><span class="badge alert-success pull-right">Success</span>' + file.name + ' ' + message + '</a>');
		});
		UploadLogic.flow.on('fileError', function(file, message){
			removeFromCookie(file.name);
		    $('#fileLog').append('<a href="#" class="list-group-item list-group-item-danger"><span class="badge alert-danger pull-right">Error</span>' + file.name + ' ' + message + '</a>');
		});

		UploadLogic.flow.on('progress', UploadLogic.onProgress);

		UploadLogic.flow.on('complete', function(){
		    $('.progress-bar').css('width', '0%');
		    $('#fileLog').append('<a href="#" class="list-group-item list-group-item-success"><span class="badge alert-success pull-right">Success</span>All upload completed</a>');
		    $('#upload_form').show();
		});
		UploadLogic.flow.on('uploadStart', function(){
		    $('#upload_form').hide();
			$('#list-of-not-uploded-files-holder').hide();
		});

		$('#js-upload-form').submit(function(ev){
			ev.preventDefault();
			
			setTimeout(function(){

				UploadLogic.flow.upload();
				
			},1);
		});

		$('a.trigger-browse-files').on("click", function(ev){
			$('input[type=file]').click();
	    	return false;
		});

	},

	onProgress: function () {
		
		if (UploadLogic.flow.progress == 0) {
			return
		}

		if (UploadLogic.firstPart == false) {
			UploadLogic.flow.pause();
			$.ajax({
				url: "/check/",
				dataType: "json"
			}).done(function(data) {
				console.log(data)
			});
		}
		else {
			var progress = UploadLogic.flow.progress() * 100;
			if (progress < 100) {
				$('.progress-bar').css('width', progress + '%');
			}
		}
	}
};


jQuery(document).ready(function(){

	UploadLogic.init();

});