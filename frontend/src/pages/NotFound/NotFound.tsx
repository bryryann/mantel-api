import styles from './NotFound.module.css';

const NotFound = () => {
    return (
        <main className={styles.container}>
            <div>
                <h2>404 Not Found.</h2>
                <p>The page you're looking for doesn't exist or has been removed.</p>
            </div>
        </main>
    );
};

export default NotFound;
