
import styled from 'styled-components';

export const Group = styled.div`
    display: flex;
    flex-direction: row;
    gap: 1em;
    width: 100%;
    height: min-content;
`;

export const Container = styled.div`
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: flex-start;
    align-items: center;
    flex-direction: column;
    overflow-y: auto;
    gap: 1em;
    padding: 1em;
    position: relative;
`;
