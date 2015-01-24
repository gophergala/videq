
/**
 * Video Background
 */
var Video = {


	init : function () {

		var root = $('body.home'),
			vid_cont = $('#intro_vid', root);

		if( vid_cont.length > 0 && !Modernizr.touch && Modernizr.mq('only screen and (min-width: 781px)') ) {


			var html = '<video autoplay loop poster="/resources/vid/loop.jpg">'+
					   '<source src="/resources/vid/loop.mp4" type="video/mp4" />'+
					   '<source src="/resources/vid/loop.webm" type="video/webm" />'+
					   '<source src="/resources/vid/loop.ogv" type="video/ogg" />'+
					   '</video>';

			vid_cont.prepend(html);
		}

	}

};


$(function(){


	Video.init();

	var notCompletedFiles = getFilesListFromCookie();
	if (notCompletedFiles != false) {
		$.each(notCompletedFiles, function(ix, fileName){
			$('#list-of-not-uploded-files').append('<li>' + fileName + '</li>');
		});
		$('#list-of-not-uploded-files-holder').show();
	}
	
	$("#drop-zone").hover(function(){

		$(".gopher").toggleClass("over");

	});


	$(".bswitch").bootstrapSwitch();

});