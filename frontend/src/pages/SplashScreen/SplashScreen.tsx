import { useState, useEffect } from 'react';
import Button from '@/components/Button';
import styles from './SplashScreen.module.css';
import logoColor from '@/assets/images/logo-color.png';
import monogramColor from '@/assets/images/monogram-color.png';

const SplashScreen = () => {
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

                <div className={styles.headerBtns}>
                    <Button variant='outline'>Sign Up</Button>
                    <Button variant='primary'>Log In</Button>
                </div>
            </header>
            <div>
            </div>
        </div>
    );
};

export default SplashScreen;
