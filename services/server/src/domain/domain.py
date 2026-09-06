from lottery.bet import Bet

def string_to_bet(bet_string, agency_id):
    fields = bet_string.split(',')
    first_name, last_name, document, birthdate, number = fields
    return Bet(int(agency_id), first_name, last_name, int(document), birthdate, int(number))

def bet_to_string(bet):
    return f"{bet.first_name},{bet.last_name},{bet.document},{bet.birthdate},{bet.number}"

def strings_to_bets(payload, agency_id):
    return [string_to_bet(line, agency_id) for line in payload]