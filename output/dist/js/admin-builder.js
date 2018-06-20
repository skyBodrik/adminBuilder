$(function() {
    let fieldsDescription = {{ JSON ._fieldsDescription }};
    console.log(fieldsDescription);
    //alert(location.href);
    //var q = Url.getQuery(location.href);
    //console.log(Url.parseQuery(q));
    $.ajax({
        url: "{{ URL ._initUrl }}",
        method: "POST",
        data: Url.getQuery(location.href),
        dataType: "json",
        success: function(data) {
            if (data['status'] === null) {
                console.log(transformToHtml(data['data']['{{ ._fieldName }}'], '{{ ._fieldName }}', fieldsDescription));
            }
            console.log(data);
        }
    });
});