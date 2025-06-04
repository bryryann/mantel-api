import { useState } from "react";

type InputTypes = 'password' | 'text' | 'email';

function useInput(type: InputTypes) {
    const [value, setValue] = useState('');

    const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setValue(e.target.value);
    };

    const clear = () => {
        setValue('');
    };

    return {
        value,
        onChange,
        clear,
        bind: {
            value,
            onChange,
            type,
        },
    };
}

export default useInput;
