import { useState, useEffect } from 'react';
import LoginForm from '@/components/auth/LoginForm';
import RegisterForm from '@/components/auth/RegisterForm';
import logoColor from '@/assets/images/logo-color.png';
import monogramColor from '@/assets/images/monogram-color.png';
import styles from './AuthorizationScreen.module.css';

const AuthorizationScreen = () => {
    const [isMobile, setIsMobile] = useState(window.innerWidth <= 768);
    const [isLogin, setIsLogin] = useState(true);

    useEffect(() => {
        const handleResize = () => {
            setIsMobile(window.innerWidth <= 768);
        };

        window.addEventListener('resize', handleResize);

        return () => {
            window.removeEventListener('resize', handleResize);
        };
    }, []);

    const toggleForm = () => { setIsLogin(prev => !prev); };

    const CurrentForm = isLogin ? LoginForm : RegisterForm;
    const toggleMessage = isLogin ? "Don't have an account yet?" : "Already have an account?";
    const toggleButton = isLogin ? "Sign up" : "Log in";

    return (
        <div className={styles.container}>
            <header>
                <img
                    src={isMobile ? monogramColor : logoColor}
                    alt='Mantel Logo'
                    className='logo'
                />

                <h3>Sign in or register a new account</h3>
            </header>
            <div>
                {<CurrentForm />}
                <div className={styles.toggle}>
                    {toggleMessage}
                    <button
                        type='button'
                        onClick={toggleForm}
                    >
                        {toggleButton}
                    </button>
                </div>
            </div>
        </div>
    );
};

export default AuthorizationScreen;
