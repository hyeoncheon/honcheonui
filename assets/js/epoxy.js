/* Mix&Fix based on belows:
 *
 * Sidebar responsive
 * Bootstrap 3.3.0 Snippet by Gab89
 * https://bootsnipp.com/snippets/Zkezz
 *
 * Orange sidebar
 * Bootstrap 3.2.0 Snippet by keshavkatwe
 * https://bootsnipp.com/snippets/3xjDn
 */

$(document).ready(function(){
	$(".push_menu").click(function(){
		$("#body-wrapper").toggleClass("active");
	});

	// initially closing (with animation :-)
	if (window.location.pathname === "/") {
		$("#body-wrapper").toggleClass("active");
	}

	$(".menu a:not('.selector')").parent().removeClass("active");
	$(".menu a:not('.selector')").each(function(index) {
		if ($(this).attr('href') == document.location.pathname) {
			$(this).parent().addClass("active");
		}
	});
});

/* vim: set ts=2 sw=2 noexpandtab: */
