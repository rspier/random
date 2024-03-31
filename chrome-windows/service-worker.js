
chrome.action.onClicked.addListener(({ e }) => {
    chrome.windows.create({
        'type': 'popup',
        'height': 300,
        'width': 100,
        'url': 'switcher.html'
    })
});
