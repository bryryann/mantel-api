import useInput from '@/hooks/useInput';
import styles from './LoginForm.module.css';

const LoginForm = () => {
    const username = useInput('text');
    const password = useInput('password');

    return (
        <form className={styles.form}>
            <input placeholder='Username' {...username.bind} />
            <input placeholder='Password' {...password.bind} />
            <button type='submit'>Log In</button>
        </form>
    );
};

export default LoginForm;
