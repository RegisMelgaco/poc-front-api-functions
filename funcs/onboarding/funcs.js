const initState  = () => ([]);

function completeOnboarding(payload, state) {
    if (!state) {
        state = initState();
    }

    state.push(payload.body.onboardingKey);

    return JSON.stringify(state);
}

function getCompletedOnboardings(payload, state) {
    if (!state) {
        return JSON.stringify(initState());
    }

    return JSON.stringify(state);
}
