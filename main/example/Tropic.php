<?php

/**
  * @adminBuilder
  *     tab{name: "Первая вкладка"},
  *     requestUrl "tropic/index.php",
  *     snippetsPath "./stories",
  *     outputPath "~/var/test.html";
  */
class Story_Admin_Controller_Tropic extends Base_Controller_Admin
{

    public function indexAction()
    {
        $this->tpl = 'Tropic/index.phtml';
        $this->view->subTpl = '_main.phtml';
        $userId = $this->p('userId', $this->USER->getId());


        $game = $this->_getGame($userId);

        $this->view->storyGame = $game;

    }

    /**
     *  Получить инфу о пользователе
     *  @adminBuilder
     *      title "Информация о пользователе",
     *      render on,
     *      in userId,
     *      out userName{snippet: "textbox", caption: "Имя пользователя"},
     *      out userPass{snippet: "textbox", caption: "Пароль"};
     *  @adminBuilder test;
     **/
    public function getUserInfoAction()
    {
        $this->norender();
        $userId = $this->p('userId', $this->USER->getId());
        return [
            'userName' => 'Игорь',
            'userPass' => '123',
        ];
    }

    /**
     *  @adminBuilder render off
     *  @adminBuilder in userId, newUserName
     *  @adminBuilder out status,
     **/
    public function setUserNameAction()
    {
        $this->norender();
        $userId = $this->p('userId', $this->USER->getId());
        $questId = $this->p('newUserName', 0);
    }
}