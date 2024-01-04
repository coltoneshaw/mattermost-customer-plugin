import React from 'react';
import styled from 'styled-components';
type Params = {
    size?: string;
}

const Icon = styled.span`
    color: 'rgba(var(--center-channel-color-rgb), 0.56)',
    '&:hover': {
        backgroundColor: 'rgba(var(--center-channel-color-rgb), 0.08)',
        color: 'rgba(var(--center-channel-color-rgb), 0.72)',
    },
    svg {
        height: ${(props: Params) => props.size};
        width: ${(props: Params) => props.size};
        fill: currentColor;
    }
`;
const SaveIcon = (
    {size = '28px'}: Params,
) => {
    return (
        <Icon
            size={size}
        >
            <svg
                xmlns='http://www.w3.org/2000/svg'
                viewBox='0 0 24 24'
                height={size}
                width={size}
                fill='currentColor'
            >
                <path d='M17 3H5C3.89 3 3 3.9 3 5V19C3 20.1 3.89 21 5 21H19C20.1 21 21 20.1 21 19V7L17 3M19 19H5V5H16.17L19 7.83V19M12 12C10.34 12 9 13.34 9 15S10.34 18 12 18 15 16.66 15 15 13.66 12 12 12M6 6H15V10H6V6Z'/>
            </svg>
        </Icon>
    );
};

export {
    SaveIcon,
};
