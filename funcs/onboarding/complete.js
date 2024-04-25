function completeOnboarding({body}, state) {
    if (!state) {
        state = initState();
    }

    state.push(body.onboardingKey);

    return JSON.stringify(state);
}
