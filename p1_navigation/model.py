import torch
import torch.nn as nn
import torch.nn.functional as F

class QNetwork(nn.Module):
    """Actor (Policy) Model."""

    def __init__(self, state_size, action_size, seed):
        """Initialize parameters and build model.
        Params
        ======
            state_size (int): Dimension of each state
            action_size (int): Dimension of each action
            seed (int): Random seed
        """
        super(QNetwork, self).__init__()
        self.seed = torch.manual_seed(seed)
        "*** YOUR CODE HERE ***"
        self.output_size = 64
        self.hidden_size = 64
        self.model = nn.Sequential(nn.Linear(state_size, self.output_size),
                                   nn.ReLU(),
                                   nn.Linear(self.output_size, self.hidden_size),
                                   nn.ReLU(),
                                   nn.Linear(self.hidden_size, action_size))

    def forward(self, state):
        """Build a network that maps state -> action values."""
        return self.model.forward(state)
