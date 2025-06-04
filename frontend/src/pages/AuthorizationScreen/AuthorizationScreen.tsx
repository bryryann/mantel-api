import { useState, useEffect } from 'react';
import LoginForm from '@/components/auth/LoginForm';
import logoColor from '@/assets/images/logo-color.png';
import monogramColor from '@/assets/images/monogram-color.png';
import styles from './AuthorizationScreen.module.css';

const AuthorizationScreen = () => {
    const [isMobile, setIsMobile] = useState(window.innerWidth <= 768);

    useEffect(() => {
        const handleResize = () => {
            setIsMobile(window.innerWidth <= 768);
        };

        window.addEventListener('resize', handleResize);

        return () => {
            window.removeEventListener('resize', handleResize);
        };
    }, []);

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
                <LoginForm />
                <div className={styles.toggle}>
                    Don't have an account yet?
                    <button
                        type='button'
                    >
                        Sign up.
                    </button>
                </div>
            </div>
        </div>
    );
};

export default AuthorizationScreen;
