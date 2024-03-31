
async function windowChange() {
    chrome.windows.getAll({ 'populate': true }, async function (ws) {

        // clear existing content...
        let div = document.getElementById("content")
        while (div.firstChild) {
            div.removeChild(div.firstChild);
        }

        let current = await chrome.windows.getCurrent()
        for (w of ws) {
            if (w.id == current.id) { // don't include this extension
                continue;
            }
            let at = activeTab(w.tabs)

            let row = document.createElement("a");
            if (w.focused) {
                row.className += " focused"
            }
            let img = document.createElement("img");
            if (at.favIconUrl != "") {
                img.src = at.favIconUrl
            }
            img.className = "favIcon"
            let p = document.createElement("span");
            p.innerText = at.title;

            row.append(img, p)

            let id = w.id
            row.addEventListener('click', () => {
                console.log("opening " + id)
                chrome.windows.update(id, { 'focused': true })
            })

            div.append(row)

        }



    });
}

chrome.windows.onCreated.addListener(windowChange);
chrome.windows.onRemoved.addListener(windowChange);
chrome.windows.onFocusChanged.addListener(windowChange);

function activeTab(tabs) {
    for (tab of tabs) {
        if (tab.active) {
            return tab;
        }
    }
}

windowChange();