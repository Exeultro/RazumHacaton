export function useValidation() {
    const rules = {
        required: (value) => !!value || 'Это поле обязательно для заполнения',
        email: (value) => {
            const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            return pattern.test(value) || 'Некорректный email';
        },
        minLength: (min) => (value) =>
            (value && value.length >= min) || `Минимум ${min} символов`,
    };

    const validateForm = (formData, formRules) => {
        const errors = {};
        let isValid = true;

        for (const key in formRules) {
            for (const rule of formRules[key]) {
                const result = rule(formData[key]);
                if (result !== true) {
                    errors[key] = result;
                    isValid = false;
                    break;
                }
            }
        }
        return { isValid, errors };
    };

    return { rules, validateForm };
}