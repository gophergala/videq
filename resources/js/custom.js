

/**
 * Video Background
 */
var Video = {


	init : function () {

		var root = $('body.home'),
			vid_cont = $('#intro_vid', root);

		if( vid_cont.length > 0 && !Modernizr.touch && Modernizr.mq('only screen and (min-width: 781px)') ) {


			var html = '<video autoplay loop poster="/resources/vid/loop.jpg">' +
					   '<source src="/resources/vid/loop.mp4" type="video/mp4" />' +
					   '<source src="/resources/vid/loop.webm" type="video/webm" />' +
					   '<source src="/resources/vid/loop.ogv" type="video/ogg" />' +
					   '</video>';

			vid_cont.prepend(html);
		}

	}

};

/**
 * Screeen switcher
 */
var Screen = {

	active : false,

	init : function (name) {

		var first = name!==undefined ? name : 'screen-drop-zone';

		Screen.show(first);
	},

	show : function (name) {

		if($(".screen." + name ).length>0)
		{
			this.active = name;

			Screen.addScreenCookie(this.active);

			$(".screen").hide();
			$(".screen." + name ).show();
		}
	
	},


	addScreenCookie : function (screenName) {
		
		Screen.removeScreenCookie();

		$.cookie('screen', screenName, { expires: 7, path: '/' });
	},


	removeScreenCookie : function () {
		$.cookie('screen', '', { expires: 7, path: '/' });
	},

	getScreentFromCookie : function () {
		var cookieStorage = $.cookie('screen');

		if (cookieStorage == undefined || cookieStorage == ""){
			return false;
		}

		return cookieStorage;
	}


};


/**
 * Flash messages
 */
var Msg = {

	init : function () {
		this.hide();
	},

	hide : function () {
		$(".flashmsg div").hide();
	},

	show : function () {
		$(".flashmsg div").show();
	},

	success : function (msg) {
		this.hide();
		$(".flashmsg .imoon").addClass("icon-circle-check");
		$(".flashmsg div").attr("class", "").addClass("success");
		$(".flashmsg div p").html(msg);
		this.show();
	},

	error : function (msg) {
		this.hide();
		$(".flashmsg .imoon").addClass("icon-notification");
		$(".flashmsg div").attr("class", "").addClass("error");
		$(".flashmsg div p").html(msg);
		this.show();
	},

	info : function (msg) {
		this.hide();
		$(".flashmsg .imoon").addClass("icon-info");
		$(".flashmsg div").attr("class", "").addClass("info");
		$(".flashmsg div p").html(msg);
		this.show();
	},

	warning : function (msg) {
		this.hide();
		$(".flashmsg .imoon").addClass("icon-warning");
		$(".flashmsg div").attr("class", "").addClass("warning");
		$(".flashmsg div p").html(msg);
		this.show();
	}

};



$(function(){


	Video.init();

	var default_screen = Screen.getScreentFromCookie();

	if(default_screen!==undefined && default_screen!="")
	{
		if(default_screen=="screen-progress-bar" || default_screen=="screen-converting-bar" || default_screen=="screen-download-bar")
		{
			if(default_screen=="screen-converting-bar")
			{
				UploadLogic.isDone();
			}			

			if(default_screen=="screen-download-bar")
			{
				UploadLogic.isDone();
			}

			Screen.init(default_screen);
		}
		else
		{
			Screen.init();
		}
	}
	else
	{
		Screen.init(); // 'screen-download-bar'
	}


	Msg.init();

	var notCompletedFiles = getFilesListFromCookie();
	if (notCompletedFiles != false) {

		var fileString = "";

		$.each(notCompletedFiles, function(ix, fileName){

			fileString += fileName;
		});
		

		Msg.warning("Drag and drop following files and click upload to resume them: " + fileString);
	}
	
	$("#drop-zone").hover(function(){

		$(".gopher").toggleClass("over");

	});


	$(".bswitch").bootstrapSwitch();

	$(".trigger_new_session").on("click", function(e){

		e.preventDefault();
		Screen.removeScreenCookie();
		$.cookie('files', '', { expires: 7, path: '/' });

		location.href = '';

		$.ajax({
			url: "/restart/",
			dataType: "json"
		}).done(function(data){
			

		});

	});

});