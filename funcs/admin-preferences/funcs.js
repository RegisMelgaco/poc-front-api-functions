const initState  = () => ({"ui":"TEMPL+HTMX rules"});

function updatePrefs(payload, state) {
    if (payload.body == null) {
        return JSON.stringify(initState());
    }

    return JSON.stringify(payload.body);
}

function getPrefs(payload, state) {
    if (!state) {
        return JSON.stringify(initState());
    }

    return JSON.stringify(state);
}
