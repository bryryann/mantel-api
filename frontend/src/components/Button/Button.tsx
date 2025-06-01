import { ButtonHTMLAttributes, ReactNode } from 'react';
import styles from './Button.module.css';

type ButtonVariant = 'primary' | 'secondary' | 'danger' | 'outline';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    children: ReactNode;
    variant?: ButtonVariant;
};

const Button: React.FC<ButtonProps> = ({
    children,
    variant = 'primary',
    className = '',
    disabled = false,
    ...props
}) => {
    const computedClass = `${styles.button} ${styles[variant]} ${className}`;

    return (
        <button
            className={computedClass}
            disabled={disabled}
            {...props}
        >
            {children}
        </button>
    );
};

export default Button;
