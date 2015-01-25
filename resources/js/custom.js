
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
			$(".screen").hide();
			$(".screen." + name ).show();
		}
	
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
	Screen.init(); // 'screen-download-bar'
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

});