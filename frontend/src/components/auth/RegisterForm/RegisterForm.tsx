import useInput from '@/hooks/useInput';
import styles from '../auth.module.css';

const RegisterForm = () => {
    const username = useInput('text');
    const email = useInput('email');
    const password = useInput('password');

    return (
        <form className={styles.form}>
            <input placeholder='Username' {...username.bind} />
            <input placeholder='Email' {...email.bind} />
            <input placeholder='Password' {...password.bind} />
            <button type='submit'>Create Account</button>
        </form>
    );
};

export default RegisterForm;

