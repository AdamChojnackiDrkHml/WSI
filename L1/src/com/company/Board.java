package com.company;



import java.util.*;

import static com.company.Constants.BOARD_SIZE;

public class Board
{
    private static Board board_singleton = null;
    private byte[]  board_content = null;

    private Board()
    {
        board_content = new byte[BOARD_SIZE*BOARD_SIZE];


        for(int i = 0; i < (BOARD_SIZE*BOARD_SIZE); i++)
        {
            board_content[i] = (byte) (i+1);

        }
        board_content[BOARD_SIZE*BOARD_SIZE-1] = 0;
        shuffleArray(board_content);

        for (int i = 0; i < BOARD_SIZE*BOARD_SIZE ; i++)
        {

            if (board_content[i] == 0)
            {
                board_content[BOARD_SIZE*BOARD_SIZE-1] ^= board_content[i];
                board_content[i] ^= board_content[BOARD_SIZE*BOARD_SIZE-1];
                board_content[BOARD_SIZE*BOARD_SIZE-1] ^= board_content[i];
                break;
            }
        }
    }

    public static Board getInstance()
    {
        if(board_singleton == null)
        {
            board_singleton = new Board();
        }

        return board_singleton;
    }

    public byte[] getBoardContent()
    {
        return board_content;
    }

    private void shuffleArray(byte[] array)
    {
        int index;
        Random random = new Random();
        for (int i = array.length - 1; i > 0; i--)
        {
            index = random.nextInt(i + 1);
            if (index != i)
            {
                array[index] ^= array[i];
                array[i] ^= array[index];
                array[index] ^= array[i];
            }
        }
    }

    public static void printBoard(byte[] board_content)
    {
        for(int i = 0; i < BOARD_SIZE; i++)
        {
            for(int j = 0; j < BOARD_SIZE; j++)
            {
                if (board_content[translate2Dto1D(i,j)] == 0)
                {
                    System.out.print(" \t");
                    continue;
                }
                System.out.print("\t" + board_content[translate2Dto1D(i,j)] + " ");
            }
            System.out.println();
        }
        System.out.println();
    }

    public static byte translate2Dto1D(int i, int j)
    {
        return (byte) (i * BOARD_SIZE + j);
    }

    public static byte[] translate1Dto2D(int i)
    {
        return new byte[]{(byte)(i % BOARD_SIZE), (byte) (i / BOARD_SIZE)};
    }
    public static boolean isSolved(byte[] board_content)
    {
        for(int i = 0; i < (BOARD_SIZE*BOARD_SIZE)-1; i++)
        {
            if(board_content[i] != i+1)
            {
                return false;
            }
        }
        return board_content[BOARD_SIZE*BOARD_SIZE-1] == 0;
    }

    public static ArrayList<byte[]> getOutcomes(byte[] board_content)
    {
        int[] zeroCords = {0,0};
        int zeroPos = 0;
        ArrayList<byte[]> neighbours = new ArrayList<>();
        for(int i = 0; i < BOARD_SIZE*BOARD_SIZE; i++)
        {
            if(board_content[i] == 0)
            {
                zeroCords[0] = i % BOARD_SIZE;
                zeroCords[1] = i / BOARD_SIZE;
                zeroPos = i;
                break;
            }
        }
        int swapPos;
        if(zeroCords[1] > 0)
        {
            swapPos = zeroPos - BOARD_SIZE;
            addNeighbourToList(zeroPos, neighbours, swapPos,board_content);

        }

        if(zeroCords[1] < BOARD_SIZE - 1)
        {
            swapPos = zeroPos + BOARD_SIZE;
            addNeighbourToList(zeroPos, neighbours, swapPos,board_content);
        }

        if(zeroCords[0] > 0)
        {
            swapPos = zeroPos - 1;
            addNeighbourToList(zeroPos, neighbours, swapPos,board_content);

        }

        if(zeroCords[0] < BOARD_SIZE - 1)
        {
            swapPos = zeroPos + 1;
            addNeighbourToList(zeroPos, neighbours, swapPos,board_content);
        }


        return neighbours;
    }

    private static void addNeighbourToList(int zeroPos, ArrayList<byte[]> neighbours, int swapPos, byte[] board_content)
    {
        byte[] outcome = new byte[BOARD_SIZE * BOARD_SIZE];

        System.arraycopy(board_content, 0, outcome, 0, BOARD_SIZE * BOARD_SIZE);
        outcome[zeroPos] = (byte) (outcome[zeroPos] ^ outcome[swapPos]);
        outcome[swapPos] = (byte) (outcome[zeroPos] ^ outcome[swapPos]);
        outcome[zeroPos] = (byte) (outcome[zeroPos] ^ outcome[swapPos]);
        neighbours.add(outcome);
    }
}
