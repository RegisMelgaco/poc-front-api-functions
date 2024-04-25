

function getCompletedOnboardings(payload, state) {
    if (!state) {
        return JSON.stringify(initState());
    }

    return JSON.stringify(state);
}