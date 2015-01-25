

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
	first: 'screen-drop-zone',

	init : function () {

		var showFirst = true;

		var cookie_screen = Screen.getScreentFromCookie();
		if(cookie_screen!==undefined && cookie_screen!="")
		{

			if(cookie_screen=="screen-progress-bar")
			{
				showFirst = false;
				Screen.show(cookie_screen);
			}

			if(cookie_screen=="screen-converting-bar")
			{
				showFirst = false;
				UploadLogic.isDone();
				Screen.show(cookie_screen);
			}

			if(cookie_screen=="screen-download-bar")
			{
				showFirst = false;
				Screen.show(cookie_screen);
				UploadLogic.atDone();
			}

			if(showFirst) Screen.show(Screen.first);

			return;
		}

		$(document).trigger("freespace");
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

	$( document ).on( "freespace", {}, function(e) {

		Screen.show(Screen.first);

		$.ajax({
			url: "/free/",
			dataType: "json"
		}).done(function(data){
			
			if(data.Procede!=undefined && data.Procede===false)
			{
				Screen.show('screen-error-bar');
				Msg.error("There was an system error. Please try again later.");
			}


		});

	});


	Video.init();
	Screen.init();
	Msg.init();


});