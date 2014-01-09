function processMessage() {
    // Parse document
    var doc = getValue("Document.Body");
    f(doc);

    // Write the results
    setValue(links);

    return 0;
}

var links = [];

var f = function(n) {
    if (n.Data == "a") {
        var index;
        for (index = 0; index < n.Attr.length; ++index) {
            var a = n.Attr[index];

            if (a.Key == "href") {
                links.push(a.Val);
                break;
            }
        }
    }


    var c;
    for (c = n.FirstChild; c !== null; c = c.NextSibling) {
        f(c);
    }
};
