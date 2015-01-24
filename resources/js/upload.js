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

		this.flow = new Flow({
			target: '/upload/',
			uploadMethod: 'POST',
			testChunks: true,
			simultaneousUploads: 1,
			prioritizeFirstAndLastChunk: true,
			progressCallbacksInterval: 0,
			singleFile: true
		});


		if (!this.flow.support) {
			alert("Browser dose not support modern upload!");
			$('#upload_form').hide();
		}
	},

	bindEvents: function () {

		this.flow.assignBrowse(document.getElementById('js-upload-files'));
		this.flow.assignDrop(document.getElementById('drop-zone'));

		this.flow.on('fileAdded', function(file, event){
			addToCookie(file.name);
			this.flow.upload();
		});

		this.flow.on('fileSuccess', function(file,message){
			removeFromCookie(file.name);
		    $('#fileLog').append('<a href="#" class="list-group-item list-group-item-success"><span class="badge alert-success pull-right">Success</span>' + file.name + ' ' + message + '</a>');
		});
		this.flow.on('fileError', function(file, message){
			removeFromCookie(file.name);
		    $('#fileLog').append('<a href="#" class="list-group-item list-group-item-danger"><span class="badge alert-danger pull-right">Error</span>' + file.name + ' ' + message + '</a>');
		});

		this.flow.on('progress', UploadLogic.onProgress);

		this.flow.on('complete', function(){
		    $('.progress-bar').css('width', '0%');
		    $('#fileLog').append('<a href="#" class="list-group-item list-group-item-success"><span class="badge alert-success pull-right">Success</span>All upload completed</a>');
		    $('#upload_form').show();
		});
		this.flow.on('uploadStart', function(){
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

	},

	onProgress: function () {
		
		if (this.flow.progress == 0) {
			return
		}

		if (UploadLogic.firstPart == false) {
			this.flow.pause();
			$.ajax({
				url: "/check/",
				dataType: "json"
			}).done(function(data) {
				console.log(data)
			});
		}
		else {
			var progress = this.flow.progress() * 100;
			if (progress < 100) {
				$('.progress-bar').css('width', progress + '%');
			}
		}
	}
};


$(function(){

	UploadLogic.init();

});