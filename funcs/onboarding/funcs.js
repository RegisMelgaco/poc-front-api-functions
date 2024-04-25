const initState  = () => ([]);

function completeOnboarding({body}, state) {
    if (!state) {
        state = initState();
    }

    state.push(body.onboardingKey);

    return JSON.stringify(state);
}

function getCompletedOnboardings(payload, state) {
    if (!state) {
        return JSON.stringify(initState());
    }

    return JSON.stringify(state);
}
