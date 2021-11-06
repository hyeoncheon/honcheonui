require("expose-loader?exposes=$,jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
$(() => {

});

$(document).ready(function(){
	var now = moment();

	$('.time').each(function(i, e) {
		var format = $(e).attr('form');
		if (format == undefined) {
			format = "YYYY-MM-DD hh:mm";
		}
		var time = moment($(e).text());
		var html = '<span title="' + time.format() + '">';
		if(now.diff(time, 'days') <= 28) {
			html += time.fromNow();
		} else {
			html += time.format(format);
		}
		$(e).html(html + '</span>');
	});
});
/* vim: set ts=2 sw=2 noexpandtab: */
