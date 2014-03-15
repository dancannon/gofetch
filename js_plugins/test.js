var links = [];

var processNode = function(n) {
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

// Parse document
processNode(document.Body);

// Write the results
setValue(links);
